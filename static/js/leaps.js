/*
Copyright (c) 2014 Ashley Jeffs

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, sub to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

/*jshint newcap: false*/

var leap_client = {};

(function() {
"use strict";

/*--------------------------------------------------------------------------------------------------
 */
/* leap_model is an object designed to keep track of the inbound and outgoing transforms
 * for a local document, and updates the caller with the appropriate actions at each stage.
 *
 * leap_model has three states:
 * 1. READY     - No pending sends, transforms received can be applied instantly to local document.
 * 2. SENDING   - Transforms are being sent and we're awaiting the corrected version of those
 *                transforms.
 * 3. BUFFERING - A corrected version has been received for our latest send but we're still waiting
 *                for the transforms that came before that send to be received before moving on.
 */
var leap_model = function(base_version) {
	this.READY = 1;
	this.SENDING = 2;
	this.BUFFERING = 3;

	this._leap_state = this.READY;

	this._corrected_version = 0;
	this._version = base_version;

	this._unapplied = [];
	this._unsent = [];
	this._sending = null;
};

/* _validate_transforms iterates an array of transform objects and validates that each transform
 * contains the correct fields. Returns an error message as a string if there was a problem.
 */
leap_model.prototype._validate_transforms = function(transforms) {
	for ( var i = 0, l = transforms.length; i < l; i++ ) {
		var tform = transforms[i];

		if ( typeof(tform.position) !== "number" ) {
			tform.position = parseInt(tform.position);
			if ( isNaN(tform.position) ) {
				return "transform contained NaN value for position: " + JSON.stringify(tform);
			}
		}
		if ( tform.num_delete !== undefined ) {
			if ( typeof(tform.num_delete) !== "number" ) {
				tform.num_delete = parseInt(tform.num_delete);
				if ( isNaN(tform.num_delete) ) {
					return "transform contained NaN value for num_delete: " + JSON.stringify(tform);
				}
			}
		} else {
			tform.num_delete = 0;
		}
		if ( tform.version !== undefined && typeof(tform.version) !== "number" ) {
			tform.version = parseInt(tform.version);
			if ( isNaN(tform.version) ) {
				return "transform contained NaN value for version: " + JSON.stringify(tform);
			}
		}
		if ( tform.insert !== undefined ) {
			if ( typeof(tform.insert) !== "string" ) {
				return "transform contained non-string value for insert: " + JSON.stringify(tform);
			}
		} else {
			tform.insert = "";
		}
	}
};

/* _validate_updates iterates an array of user update objects and validates that each update
 * contains the correct fields. Returns an error message as a string if there was a problem.
 */
leap_model.prototype._validate_updates = function(user_updates) {
	for ( var i = 0, l = user_updates.length; i < l; i++ ) {
		var update = user_updates[i];

		if ( undefined !== update.position &&
		    "number" !== typeof(update.position) ) {
			update.position = parseInt(update.position);
			if ( isNaN(update.position) ) {
				return "update contained NaN value for position: " + JSON.stringify(update);
			}
		}
		//the following causes error in mac, comment it out
		//console.log("Type of update message: " + typeof(update.message) );
		//if ( undefined !== update.message &&
		//    "string" !== typeof(update.message) ) {
		//	return "update contained invalid type for message: " + JSON.stringify(update);
		//}
		if ( undefined !== update.active &&
		    "boolean" !== typeof(update.active) ) {
			if ("string" !== typeof(update.active)) {
				return "update contained invalid type for active: " + JSON.stringify(update);
			}
			update.active = ("true" === update.active);
		}
		//the following causes error in mac, comment it out
		//if ( undefined === update.user_id ||
		//    "string" !== typeof(update.user_id) ) {
		//	return "update contained invalid type for user_id: " + JSON.stringify(update);
		//}
	}
};

/* merge_transforms takes two transforms (the next to be sent, and the one that follows) and
 * attempts to merge them into one transform. This will not be possible with some combinations, and
 * the function returns a boolean to indicate whether the merge was successful.
 */
leap_model.prototype._merge_transforms = function(first, second) {
	var overlap, remainder;

	if ( first.position + first.insert.length === second.position ) {
		first.insert = first.insert + second.insert;
		first.num_delete += second.num_delete;
		return true;
	}
	if ( second.position === first.position ) {
		remainder = Math.max(0, second.num_delete - first.insert.length);
		first.num_delete += remainder;
		first.insert = second.insert + first.insert.slice(second.num_delete);
		return true;
	}
	if ( second.position > first.position && second.position < ( first.position + first.insert.length ) ) {
		overlap = second.position - first.position;
		remainder = Math.max(0, second.num_delete - (first.insert.length - overlap));
		first.num_delete += remainder;
		first.insert = first.insert.slice(0, overlap) + second.insert + first.insert.slice(overlap + second.num_delete);
		return true;
	}
	return false;
};

/* collide_transforms takes an unapplied transform from the server, and an unsent transform from the
 * client and modifies both transforms.
 *
 * The unapplied transform is fixed so that when applied to the local document is unaffected by the
 * unsent transform that has already been applied. The unsent transform is fixed so that it is
 * unaffected by the unapplied transform when submitted to the server.
 */
leap_model.prototype._collide_transforms = function(unapplied, unsent) {
	var earlier, later;

	if ( unapplied.position <= unsent.position ) {
		earlier = unapplied;
		later = unsent;
	} else {
		earlier = unsent;
		later = unapplied;
	}
	if ( earlier.num_delete === 0 ) {
		later.position += earlier.insert.length;
	} else if ( ( earlier.num_delete + earlier.position ) <= later.position ) {
		later.position += ( earlier.insert.length - earlier.num_delete );
	} else {
		var pos_gap = later.position - earlier.position;
		var over_hang = Math.min(later.insert.length, earlier.num_delete - pos_gap);
		var excess = Math.max(0, (earlier.num_delete - pos_gap));

		// earlier changes
		if ( excess > later.num_delete ) {
			earlier.num_delete += later.insert.length - later.num_delete;
			earlier.insert = earlier.insert + later.insert;
		} else {
			earlier.num_delete = pos_gap;
		}
		// later changes
		later.num_delete = Math.min(0, later.num_delete - excess);
		later.position = earlier.position + earlier.insert.length;
	}
};

/*--------------------------------------------------------------------------------------------------
 */

/* _resolve_state will prompt the leap_model to re-evalutate its current state for validity. If this
 * state is determined to no longer be appropriate then it will return an object containing the
 * following actions to be performed.
 */
leap_model.prototype._resolve_state = function() {
	switch (this._leap_state) {
	case this.READY:
	case this.SENDING:
		return {};
	case this.BUFFERING:
		if ( ( this._version + this._unapplied.length ) >= (this._corrected_version - 1) ) {

			this._version += this._unapplied.length + 1;
			var to_collide = [ this._sending ].concat(this._unsent);
			var unapplied = this._unapplied;

			this._unapplied = [];

			for ( var i = 0, li = unapplied.length; i < li; i++ ) {
				for ( var j = 0, lj = to_collide.length; j < lj; j++ ) {
					this._collide_transforms(unapplied[i], to_collide[j]);
				}
			}

			this._sending = null;

			if ( this._unsent.length > 0 ) {
				this._sending = this._unsent.shift();
				while ( this._unsent.length > 0 && this._merge_transforms(this._sending, this._unsent[0]) ) {
					this._unsent.shift();
				}
				this._sending.version = this._version + 1;
				this._leap_state = this.SENDING;
				return { send : this._sending, apply : unapplied };
			} else {
				this._leap_state = this.READY;
				return { apply : unapplied };
			}
		}
	}
	return {};
};

/* correct is the function to call following a "correction" from the server, this correction value
 * gives the model the information it needs to determine which changes are missing from our model
 * from before our submission was accepted.
 */
leap_model.prototype.correct = function(version) {
	switch (this._leap_state) {
	case this.READY:
	case this.BUFFERING:
		return { error : "received unexpected correct action" };
	case this.SENDING:
		this._leap_state = this.BUFFERING;
		this._corrected_version = version;
		return this._resolve_state();
	}
	return {};
};

/* submit is the function to call when we wish to submit more local changes to the server. The model
 * will determine whether it is currently safe to dispatch those changes to the server, and will
 * also provide each change with the correct version number.
 */
leap_model.prototype.submit = function(transform) {
	switch (this._leap_state) {
	case this.READY:
		this._leap_state = this.SENDING;
		transform.version = this._version + 1;
		this._sending = transform;
		return { send : transform };
	case this.BUFFERING:
	case this.SENDING:
		this._unsent = this._unsent.concat(transform);
	}
	return {};
};

/* receive is the function to call when we have received transforms from our server. If we have
 * recently dispatched transforms and have yet to receive our correction then it is unsafe to apply
 * these changes to our local document, so the model will keep return these transforms to us when it
 * is known to be safe.
 */
leap_model.prototype.receive = function(transforms) {
	var expected_version = this._version + this._unapplied.length + 1;
	if ( (transforms.length > 0) && (transforms[0].version !== expected_version) ) {
		return { error :
			("Received unexpected transform version: " + transforms[0].version +
				", expected: " + expected_version) };
	}

	switch (this._leap_state) {
	case this.READY:
		this._version += transforms.length;
		return { apply : transforms };
	case this.BUFFERING:
		this._unapplied = this._unapplied.concat(transforms);
		return this._resolve_state();
	case this.SENDING:
		this._unapplied = this._unapplied.concat(transforms);
	}
	return {};
};

/*--------------------------------------------------------------------------------------------------
 */

/* leap_client is the main tool provided to allow an easy and stable interface for connecting to a
 * leaps server.
 */
leap_client = function() {
	this._socket = null;
	this._document_id = null;

	this._model = null;

	this._cursor_position = 0;

	this.EVENT_TYPE = {
		CONNECT: "connect",
		DISCONNECT: "disconnect",
		DOCUMENT: "document",
		TRANSFORMS: "transforms",
		USER: "user",
		ERROR: "error"
	};

	// Milliseconds period between cursor position updates to server
	this._POSITION_POLL_PERIOD = 500;

	this._events = {};
};

/* subscribe_event, attach a function to an event of the leap_client. Use this to subscribe to
 * transforms, document responses and errors etc. Returns a string if an error occurrs.
 */
leap_client.prototype.subscribe_event = function(name, subscriber) {
	if ( typeof(subscriber) !== "function" ) {
		return "subscriber was not a function";
	}
	var targets = this._events[name];
	if ( targets !== undefined && targets instanceof Array ) {
		targets.push(subscriber);
	} else {
		this._events[name] = [ subscriber ];
	}
};

/* on - an alias for subscribe_event.
 */
leap_client.prototype.on = leap_client.prototype.subscribe_event;

/* clear_subscribers, removes all functions subscribed to an event.
 */
leap_client.prototype.clear_subscribers = function(name) {
	this._events[name] = [];
};

/* dispatch_event, sends args to all subscribers of an event.
 */
leap_client.prototype._dispatch_event = function(name, args) {
	var targets = this._events[name];
	if ( targets !== undefined && targets instanceof Array ) {
		for ( var i = 0, l = targets.length; i < l; i++ ) {
			if (typeof(targets[i]) === "function") {
				targets[i].apply(this, args);
			}
		}
	}
};

/* _do_action is a call that acts accordingly provided an action_obj from our leap_model.
 */
leap_client.prototype._do_action = function(action_obj) {
	if ( action_obj.error !== undefined ) {
		return action_obj.error;
	}
	if ( action_obj.apply !== undefined && action_obj.apply instanceof Array ) {
		this._dispatch_event(this.EVENT_TYPE.TRANSFORMS, [ action_obj.apply ]);
	}
	if ( action_obj.send !== undefined && action_obj.send instanceof Object ) {
		this._socket.send(JSON.stringify({
			command : "submit",
			transform : action_obj.send
		}));
	}
};

/* _process_message is a call that takes a server provided message object and decides the
 * appropriate action to take. If an error occurs during this process then an error message is
 * returned.
 */
leap_client.prototype._process_message = function(message) {
	var validate_error, action_obj, action_err;

	if ( message.response_type === undefined || typeof(message.response_type) !== "string" ) {
		return "message received did not contain a valid type";
	}

	switch (message.response_type) {
	case "document":
		if ( null === message.leap_document ||
		   "object" !== typeof(message.leap_document) ||
		   "string" !== typeof(message.leap_document.id) ||
		   "string" !== typeof(message.leap_document.content) ) {
			return "message document type contained invalid document object";
		}
		if ( message.version <= 0 ) {
			return "message document received but without valid version";
		}
		if ( this._document_id !== null && this._document_id !== message.leap_document.id ) {
			return "received unexpected document, id was mismatched: " +
				this._document_id + " != " + message.leap_document.id;
		}
		this.document_id = message.leap_document.id;
		this._model = new leap_model(message.version);
		this._dispatch_event(this.EVENT_TYPE.DOCUMENT, [ message.leap_document ]);
		break;
	case "transforms":
		if ( this._model === null ) {
			return "transforms were received before initialization";
		}
		if ( !(message.transforms instanceof Array) ) {
			return "received non array transforms";
		}
		validate_error = this._model._validate_transforms(message.transforms);
		if ( validate_error !== undefined ) {
			return "received transforms with error: " + validate_error;
		}
		action_obj = this._model.receive(message.transforms);
		action_err = this._do_action(action_obj);
		if ( action_err !== undefined ) {
			return "failed to receive transforms: " + action_err;
		}
		break;
	case "update":
		if ( null === message.user_updates ||
		   !(message.user_updates instanceof Array) ) {
			return "message update type contained invalid user_updates";
		}
		//console.log("The received message is: " + JSON.stringify(message, null, 2));
		validate_error = this._model._validate_updates(message.user_updates);
		if ( validate_error !== undefined ) {
			return "received updatess with error: " + validate_error;
		}
		for ( var i = 0, l = message.user_updates.length; i < l; i++ ) {
			this._dispatch_event(this.EVENT_TYPE.USER, [ message.user_updates[i] ]);
		}
		break;
	case "correction":
		if ( this._model === null ) {
			return "correction was received before initialization";
		}
		if ( typeof(message.version) !== "number" ) {
			message.version = parseInt(message.version);
			if ( isNaN(message.version) ) {
				return "correction received was NaN";
			}
		}
		action_obj = this._model.correct(message.version);
		action_err = this._do_action(action_obj);
		if ( action_err !== undefined ) {
			return "model failed to correct: " + action_err;
		}
		break;
	case "error":
		if ( this._socket !== null ) {
			this._socket.close();
		}
		if ( typeof(message.error) === "string" ) {
			return message.error;
		}
		return "server sent undeterminable error";
	default:
		return "message received was not a recognised type";
	}
};

/* send_transform is the function to call to send a transform off to the server. To keep the local
 * document responsive this transform should be applied to the document straight away. The
 * leap_client will decide when it is appropriate to dispatch the transform, and will manage
 * internally how incoming messages should be altered to account for the fact that the local
 * change was made out of order.
 */
leap_client.prototype.send_transform = function(transform) {
	if ( this._model === null ) {
		return "leap_client must be initialized and joined to a document before submitting transforms";
	}

	var validate_error = this._model._validate_transforms([ transform ]);
	if ( validate_error !== undefined ) {
		return validate_error;
	}

	var action_obj = this._model.submit(transform);
	var action_err = this._do_action(action_obj);
	if ( action_err !== undefined ) {
		return "model failed to submit: " + action_err;
	}
};

/* send_message - send a text message out to all other users connected to your shared document.
 */
leap_client.prototype.send_message = function(message) {
	if ( "string" !== typeof(message) ) {
		return "must supply message as a valid string value";
	}

	this._socket.send(JSON.stringify({
		command:  "update",
		message: message,
		position: this._cursor_position
	}));
};

/* update_cursor is the function to call to send the server (and all other clients) an update to your
 * current cursor position in the document, this shows others where your point of interest is in the
 * shared document.
 */
leap_client.prototype.update_cursor = function(position) {
	if ( "number" !== typeof(position) ) {
		return "must supply position as a valid integer value";
	}

	this._cursor_position = position;
	this._socket.send(JSON.stringify({
		command:  "update",
		position: this._cursor_position
	}));
};

/* join_document prompts the client to request to join a document from the server. It will return an
 * error message if there is a problem with the request.
 */
leap_client.prototype.join_document = function(id, token) {
	if ( this._socket === null || this._socket.readyState !== 1 ) {
		return "leap_client is not currently connected";
	}

	if ( typeof(id) !== "string" ) {
		return "document id was not a string type";
	}

	if ( this._document_id !== null ) {
		return "a leap_client can only join a single document";
	}

	this._document_id = id;

	this._socket.send(JSON.stringify({
		command : "find",
		token : token,
		document_id : this._document_id
	}));
};

/* create_document submits content to be created into a fresh document and then binds to that
 * document.
 */
leap_client.prototype.create_document = function(content, token) {
	if ( this._socket === null || this._socket.readyState !== 1 ) {
		return "leap_client is not currently connected";
	}

	if ( typeof(content) !== "string" ) {
		return "new document requires valid content (can be empty)";
	}

	if ( this._document_id !== null ) {
		return "a leap_client can only join a single document";
	}

	this._socket.send(JSON.stringify({
		command : "create",
		token : token,
		leap_document : {
			content : content
		}
	}));
};

/* connect is the first interaction that should occur with the leap_client after defining your event
 * bindings. This function will generate a websocket connection with the server, ready to bind to a
 * document.
 */
leap_client.prototype.connect = function(address, _websocket) {
	try {
		if ( _websocket !== undefined ) {
				this._socket = _websocket;
		} else if ( window.WebSocket !== undefined ) {
				this._socket = new WebSocket(address);
		} else {
			return "no websocket support in this browser";
		}
	} catch(e) {
		return "socket connection failed: " + e.message;
	}

	var leap_obj = this;

	this._socket.onmessage = function(message) {
		var message_text = message.data;
		var message_obj;

		try {
			message_obj = JSON.parse(message_text);
		} catch (e) {
			leap_obj._dispatch_event.apply(leap_obj,
				[ leap_obj.EVENT_TYPE.ERROR,
					[ JSON.stringify(e.message) + " (" + e.lineNumber + "): " + message_text ] ]);
			return;
		}

		var err = leap_obj._process_message.apply(leap_obj, [ message_obj ]);
		if ( typeof(err) === "string" ) {
			leap_obj._dispatch_event.apply(leap_obj, [ leap_obj.EVENT_TYPE.ERROR, [ err ] ]);
		}
	};

	this._socket.onclose = function() {
		if ( undefined !== leap_obj._heartbeat ) {
			clearTimeout(leap_obj._heartbeat);
		}
		leap_obj._dispatch_event.apply(leap_obj, [ leap_obj.EVENT_TYPE.DISCONNECT, [] ]);
	};

	this._socket.onopen = function() {
		leap_obj._heartbeat = setInterval(function() {
			leap_obj._socket.send(JSON.stringify({
				command : "ping"
			}));
		}, 5000); // MAGIC NUMBER OH GOD, we should have a config object.
		leap_obj._dispatch_event.apply(leap_obj, [ leap_obj.EVENT_TYPE.CONNECT, arguments ]);
	};

	this._socket.onerror = function() {
		if ( undefined !== leap_obj._heartbeat ) {
			clearTimeout(leap_obj._heartbeat);
		}
		leap_obj._dispatch_event.apply(leap_obj, [ leap_obj.EVENT_TYPE.ERROR, [ "socket connection error" ] ]);
	};
};

/* Close the connection to the document and halt all operations.
 */
leap_client.prototype.close = function() {
	if ( undefined !== this._heartbeat ) {
		clearTimeout(this._heartbeat);
	}
	if ( this._socket !== null && this._socket.readyState === 1 ) {
		this._socket.close();
		this._socket = null;
	}
	this.document_id = undefined;
	this._model = undefined;
};

/*--------------------------------------------------------------------------------------------------
 */

/* leap_apply is a function that applies a single transform to content and returns the result.
 */
var leap_apply = function(transform, content) {
	var num_delete = 0, to_insert = "";

	if ( typeof(transform.position) !== "number" ) {
		return content;
	}

	if ( typeof(transform.num_delete) === "number" ) {
		num_delete = transform.num_delete;
	}

	if ( typeof(transform.insert) === "string" ) {
		to_insert = transform.insert;
	}

	var first = content.slice(0, transform.position);
	var second = content.slice(transform.position + num_delete, content.length);
	return first + to_insert + second;
};

leap_client.prototype.apply = leap_apply;

/*--------------------------------------------------------------------------------------------------
 */

try {
	if ( module !== undefined && typeof(module) === "object" ) {
		module.exports = {
			client : leap_client,
			apply : leap_apply,
			_model : leap_model
		};
	}
} catch(e) {
}

/*--------------------------------------------------------------------------------------------------
 */

})();
/*
Copyright (c) 2014 Ashley Jeffs

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, sub to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

/*jshint newcap: false*/

(function() {
"use strict";

/*--------------------------------------------------------------------------------------------------
 */

/*
_create_leaps_ace_marker - creates a marker for displaying the cursor positions of other users in an
ace editor.
*/
var _create_leaps_ace_marker = function(ace_editor) {
	var marker = {};

	marker.draw_handler = null;
	marker.clear_handler = null;
	marker.cursors = [];

	marker.update = function(html, markerLayer, session, config) {
		if ( typeof marker.clear_handler === 'function' ) {
			marker.clear_handler();
		}
		var cursors = marker.cursors;
		for (var i = 0; i < cursors.length; i++) {
			var pos = cursors[i].position;
			var screenPos = session.documentToScreenPosition(pos);

			var height = config.lineHeight;
			var width = config.characterWidth;
			var top = markerLayer.$getTop(screenPos.row, config);
			var left = markerLayer.$padding + screenPos.column * width;

			var stretch = 4;

			if ( typeof marker.draw_handler === 'function' ) {
				var content = (marker.draw_handler(
					cursors[i].user_id, height, top, left, screenPos.row, screenPos.column
				) || '') + '';
				html.push(content);
			} else {
				html.push(
					"<div class='LeapsAceCursor' style='",
					"height:", (height + stretch), "px;",
					"top:", (top - (stretch/2)), "px;",
					"left:", left, "px; width:", width, "px'></div>");
			}
		}
	};

	marker.redraw = function() {
		marker.session._signal("changeFrontMarker");
	};

	marker.updateCursor = function(user) {
		var cursors = marker.cursors, current, i, l;
		for ( i = 0, l = cursors.length; i < l; i++ ) {
			if ( cursors[i].user_id === user.user_id ) {
				current = cursors[i];
				current.position = marker.session.getDocument().indexToPosition(user.position, 0);
				current.updated = new Date().getTime();
				break;
			}
		}
		if ( undefined === current ) {
			if ( user.active ) {
				current = {
					user_id: user.user_id,
					position: marker.session.getDocument().indexToPosition(user.position, 0),
					updated: new Date().getTime()
				};
				cursors.push(current);
			}
		} else if ( !user.active ) {
			cursors.splice(i, 1);
		}

		marker.redraw();
	};

	marker.session = ace_editor.getSession();
	marker.session.addDynamicMarker(marker, true);

	return marker;
};

/* leap_bind_ace_editor takes an existing leap_client and uses it to convert an Ace web editor
 * (http://ace.c9.io) into a live leaps shared editor.
 */
var leap_bind_ace_editor = function(leap_client, ace_editor) {
	if ( null === document.getElementById("leaps-ace-style") ) {
		var node = document.createElement('style');
		node.id = "leaps-ace-style";
		node.innerHTML =
		".LeapsAceCursor {" +
			"position: absolute;" +
			"border-left: 3px solid #D11956;" +
		"}";
		document.body.appendChild(node);
	}

	this._ace = ace_editor;
	this._leap_client = leap_client;

	this._content = "";
	this._ready = false;
	this._blind_eye_turned = false;

	this._ace.setReadOnly(true);

	this._marker = _create_leaps_ace_marker(this._ace);

	var binder = this;

	this._ace.getSession().on('change', function(e) {
		binder._convert_to_transform.apply(binder, [ e ]);
	});

	this._leap_client.subscribe_event("document", function(doc) {
		binder._content = doc.content;

		binder._blind_eye_turned = true;
		binder._ace.setValue(doc.content);
		binder._ace.setReadOnly(false);
		binder._ace.clearSelection();

		var old_undo = binder._ace.getSession().getUndoManager();
		old_undo.reset();
		binder._ace.getSession().setUndoManager(old_undo);

		binder._ready = true;
		binder._blind_eye_turned = false;

		binder._pos_interval = setInterval(function() {
			var session = binder._ace.getSession(), doc = session.getDocument();
			var position = session.getSelection().getCursor();
			var index = doc.positionToIndex(position, 0);

			binder._leap_client.update_cursor.apply(binder._leap_client, [ index ]);
		}, leap_client._POSITION_POLL_PERIOD);
	});

	this._leap_client.subscribe_event("transforms", function(transforms) {
		for ( var i = 0, l = transforms.length; i < l; i++ ) {
			binder._apply_transform.apply(binder, [ transforms[i] ]);
		}
	});

	this._leap_client.subscribe_event("disconnect", function() {
		binder._ace.setReadOnly(true);

		if ( undefined !== binder._pos_interval ) {
			clearTimeout(binder._pos_interval);
		}
	});

	this._leap_client.subscribe_event("user", function(user) {
		binder._marker.updateCursor.apply(binder._marker, [ user ]);
	});

	this._leap_client.ACE_set_cursor_handler = function(handler, clear_handler) {
		binder.set_cursor_handler(handler, clear_handler);
	};
};

/* set_cursor_handler, sets the method call that returns a cursor marker. Also adds an optional
 * clear_handler which is called before each individual cursor is drawn (use it to clear all outside
 * markers before redrawing).
 */
leap_bind_ace_editor.prototype.set_cursor_handler = function(handler, clear_handler) {
	if ( 'function' === typeof handler ) {
		this._marker.draw_handler = handler;
	}
	if ( 'function' === typeof clear_handler ) {
		this._marker.clear_handler = clear_handler;
	}
};

/* apply_transform, applies a single transform to the ace document.
 */
leap_bind_ace_editor.prototype._apply_transform = function(transform) {
	this._blind_eye_turned = true;

	var edit_session = this._ace.getSession();
	var live_document = edit_session.getDocument();

	var position = live_document.indexToPosition(transform.position, 0);

	if ( transform.num_delete > 0 ) {
		edit_session.remove({
			start: position,
			end: live_document.indexToPosition(transform.position + transform.num_delete, 0)
		});
	}
	if ( typeof(transform.insert) === "string" && transform.insert.length > 0 ) {
		edit_session.insert(position, transform.insert);
	}

	this._blind_eye_turned = false;

	this._content = this._leap_client.apply(transform, this._content);

	setTimeout((function() {
		if ( this._content !== this._ace.getValue() ) {
			this._leap_client._dispatch_event.apply(this._leap_client,
				[ this._leap_client.EVENT_TYPE.ERROR, [
					"Local editor has lost synchronization with server"
				] ]);
		}
	}).bind(this), 0);
};

/* convert_to_transform, takes an ace editor event, converts it into a transform and sends it.
 */
leap_bind_ace_editor.prototype._convert_to_transform = function(e) {
	if ( this._blind_eye_turned ) {
		return;
	}

	var tform = {};

	var live_document = this._ace.getSession().getDocument();
	var nl = live_document.getNewLineCharacter();

	switch (e.data.action) {
	case "insertText":
		tform.position = live_document.positionToIndex(e.data.range.start, 0);
		tform.insert = e.data.text;
		break;
	case "insertLines":
		tform.position = live_document.positionToIndex(e.data.range.start, 0);
		tform.insert = e.data.lines.join(nl) + nl;
		break;
	case "removeText":
		tform.position = live_document.positionToIndex(e.data.range.start, 0);
		tform.num_delete = e.data.text.length;
		break;
	case "removeLines":
		tform.position = live_document.positionToIndex(e.data.range.start, 0);
		tform.num_delete = e.data.lines.join(nl).length + nl.length;
		break;
	}

	if ( tform.insert === undefined && tform.num_delete === undefined ) {
		this._leap_client._dispatch_event.apply(this._leap_client,
			[ this._leap_client.EVENT_TYPE.ERROR, [
				"Local change resulted in invalid transform"
			] ]);
	}

	this._content = this._leap_client.apply(tform, this._content);
	var err = this._leap_client.send_transform(tform);
	if ( err !== undefined ) {
		this._leap_client._dispatch_event.apply(this._leap_client,
			[ this._leap_client.EVENT_TYPE.ERROR, [
				"Local change resulted in invalid transform: " + err
			] ]);
	}

	setTimeout((function() {
		if ( this._content !== this._ace.getValue() ) {
			this._leap_client._dispatch_event.apply(this._leap_client,
				[ this._leap_client.EVENT_TYPE.ERROR, [
					"Local editor has lost synchronization with server"
				] ]);
		}
	}).bind(this), 0);
};

/*--------------------------------------------------------------------------------------------------
 */

try {
	if ( window.leap_client !== undefined && typeof(window.leap_client) === "function" ) {
		window.leap_client.prototype.bind_ace_editor = function(ace_editor) {
			this._ace_editor = new leap_bind_ace_editor(this, ace_editor);
		};
	}
} catch (e) {
	console.error(e);
}

/*--------------------------------------------------------------------------------------------------
 */

})();
/*
Copyright (c) 2014 Ashley Jeffs

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, sub to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

/*jshint newcap: false*/

(function() {
"use strict";

/*--------------------------------------------------------------------------------------------------
 */

/* leap_bind_codemirror takes an existing leap_client and uses it to convert a codemirro web editor
 * (http://codemirror.net/) into a live leaps shared editor.
 */
var leap_bind_codemirror = function(leap_client, codemirror_object) {
	this._codemirror = codemirror_object;
	this._leap_client = leap_client;

	this._content = "";
	this._ready = false;
	this._blind_eye_turned = false;

	var binder = this;

	this._codemirror.on('beforeChange', function(instance, e) {
		binder._convert_to_transform.apply(binder, [ e ]);
	});

	this._leap_client.subscribe_event("document", function(doc) {
		binder._content = doc.content;

		binder._blind_eye_turned = true;
		binder._codemirror.getDoc().setValue(doc.content);

		binder._ready = true;
		binder._blind_eye_turned = false;

		binder._pos_interval = setInterval(function() {
			var live_document = binder._codemirror.getDoc();
			var position = live_document.indexFromPos(live_document.getCursor());
			binder._leap_client.update_cursor.apply(binder._leap_client, [ position ]);
		}, leap_client._POSITION_POLL_PERIOD);
	});

	this._leap_client.subscribe_event("transforms", function(transforms) {
		for ( var i = 0, l = transforms.length; i < l; i++ ) {
			binder._apply_transform.apply(binder, [ transforms[i] ]);
		}
	});

	this._leap_client.subscribe_event("disconnect", function() {
		if ( undefined !== binder._pos_interval ) {
			clearTimeout(binder._pos_interval);
		}
	});
};

/* apply_transform, applies a single transform to the codemirror document
 */
leap_bind_codemirror.prototype._apply_transform = function(transform) {
	this._blind_eye_turned = true;

	var live_document = this._codemirror.getDoc();
	var start_position = live_document.posFromIndex(transform.position), end_position = start_position;

	if ( transform.num_delete > 0 ) {
		end_position = live_document.posFromIndex(transform.position + transform.num_delete);
	}

	var insert = "";
	if ( typeof(transform.insert) === "string" && transform.insert.length > 0 ) {
		insert = transform.insert;
	}

	live_document.replaceRange(insert, start_position, end_position);

	this._blind_eye_turned = false;

	this._content = this._leap_client.apply(transform, this._content);

	setTimeout((function() {
		if ( this._content !== this._codemirror.getDoc().getValue() ) {
			this._leap_client._dispatch_event.apply(this._leap_client,
				[ this._leap_client.EVENT_TYPE.ERROR, [
					"Local editor has lost synchronization with server"
				] ]);
		}
	}).bind(this), 0);
};

/* convert_to_transform, takes a codemirror edit event, converts it into a transform and sends it.
 */
leap_bind_codemirror.prototype._convert_to_transform = function(e) {
	if ( this._blind_eye_turned ) {
		return;
	}

	var tform = {};

	var live_document = this._codemirror.getDoc();
	var start_index = live_document.indexFromPos(e.from), end_index = live_document.indexFromPos(e.to);

	tform.position = start_index;
	tform.insert = e.text.join('\n') || "";

	tform.num_delete = end_index - start_index;

	if ( tform.insert.length <= 0 && tform.num_delete <= 0 ) {
		this._leap_client._dispatch_event.apply(this._leap_client,
			[ this._leap_client.EVENT_TYPE.ERROR, [
				"Change resulted in invalid transform"
			] ]);
	}

	this._content = this._leap_client.apply(tform, this._content);
	var err = this._leap_client.send_transform(tform);
	if ( err !== undefined ) {
		this._leap_client._dispatch_event.apply(this._leap_client,
			[ this._leap_client.EVENT_TYPE.ERROR, [
				"Change resulted in invalid transform: " + err
			] ]);
	}

	setTimeout((function() {
		if ( this._content !== this._codemirror.getDoc().getValue() ) {
			this._leap_client._dispatch_event.apply(this._leap_client,
				[ this._leap_client.EVENT_TYPE.ERROR, [
					"Local editor has lost synchronization with server"
				] ]);
		}
	}).bind(this), 0);
};

/*--------------------------------------------------------------------------------------------------
 */

try {
	if ( window.leap_client !== undefined && typeof(window.leap_client) === "function" ) {
		window.leap_client.prototype.bind_codemirror = function(codemirror_object) {
			this._codemirror = new leap_bind_codemirror(this, codemirror_object);
		};
	}
} catch (e) {
}

/*--------------------------------------------------------------------------------------------------
 */

})();
/*
Copyright (c) 2014 Ashley Jeffs

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, sub to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

/*jshint newcap: false*/

(function() {
"use strict";

/*--------------------------------------------------------------------------------------------------
 */

/* leap_bind_textarea takes an existing leap_client and uses it to wrap a textarea into an
 * interactive editor for the leaps document the client connects to. Returns the bound object, and
 * places any errors in the obj.error field to be checked after construction.
 */
var leap_bind_textarea = function(leap_client, text_area) {
	this._text_area = text_area;
	this._leap_client = leap_client;

	this._content = "";
	this._ready = false;
	this._text_area.disabled = true;

	var binder = this;

	if ( undefined !== text_area.addEventListener ) {
		text_area.addEventListener('input', function() {
			binder._trigger_diff();
		}, false);
	} else if ( undefined !== text_area.attachEvent ) {
		text_area.attachEvent('onpropertychange', function() {
			binder._trigger_diff();
		});
	} else {
		this.error = "event listeners not implemented on this browser, are you from the past?";
	}

	this._leap_client.subscribe_event("document", function(doc) {
		binder._content = binder._text_area.value = doc.content;
		binder._ready = true;
		binder._text_area.disabled = false;

		binder._pos_interval = setInterval(function() {
			binder._leap_client.update_cursor.apply(binder._leap_client, [ binder._text_area.selectionStart ]);
		}, leap_client._POSITION_POLL_PERIOD);
	});

	this._leap_client.subscribe_event("transforms", function(transforms) {
		for ( var i = 0, l = transforms.length; i < l; i++ ) {
			binder._apply_transform.apply(binder, [ transforms[i] ]);
		}
	});

	this._leap_client.subscribe_event("disconnect", function() {
		binder._text_area.disabled = true;
		if ( undefined !== binder._pos_interval ) {
			clearTimeout(binder._pos_interval);
		}
	});

	this._leap_client.subscribe_event("user", function(user) {
		console.log("User update: " + JSON.stringify(user));
	});
};

/* apply_transform, applies a single transform to the textarea. Also attempts to retain the original
 * cursor position.
 */
leap_bind_textarea.prototype._apply_transform = function(transform) {
	var cursor_pos = this._text_area.selectionStart;
	var cursor_pos_end = this._text_area.selectionEnd;
	var content = this._text_area.value;

	if ( transform.position <= cursor_pos ) {
		cursor_pos += (transform.insert.length - transform.num_delete);
		cursor_pos_end += (transform.insert.length - transform.num_delete);
	}

	this._content = this._text_area.value = this._leap_client.apply(transform, content);
	this._text_area.selectionStart = cursor_pos;
	this._text_area.selectionEnd = cursor_pos_end;
};

/* trigger_diff triggers whenever a change may have occurred to the wrapped textarea element, and
 * compares the old content with the new content. If a change has indeed occurred then a transform
 * is generated from the comparison and dispatched via the leap_client.
 */
leap_bind_textarea.prototype._trigger_diff = function() {
	var new_content = this._text_area.value;
	if ( !(this._ready) || new_content === this._content ) {
		return;
	}

	var i = 0, j = 0;
	while (new_content[i] === this._content[i]) {
		i++;
	}
	while ((new_content[(new_content.length - 1 - j)] === this._content[(this._content.length - 1 - j)]) &&
			((i + j) < new_content.length) && ((i + j) < this._content.length)) {
		j++;
	}

	var tform = { position : i };

	if (this._content.length !== (i + j)) {
		tform.num_delete = (this._content.length - (i + j));
	}
	if (new_content.length !== (i + j)) {
		tform.insert = new_content.slice(i, new_content.length - j);
	}

	this._content = new_content;
	if ( tform.insert !== undefined || tform.num_delete !== undefined ) {
		var err = this._leap_client.send_transform(tform);
		if ( err !== undefined ) {
			this._leap_client._dispatch_event.apply(this._leap_client,
				[ this._leap_client.EVENT_TYPE.ERROR, [
					"Local change resulted in invalid transform"
				] ]);
		}
	}
};

/*--------------------------------------------------------------------------------------------------
 */

try {
	if ( window.leap_client !== undefined && typeof(window.leap_client) === "function" ) {
		window.leap_client.prototype.bind_textarea = function(text_area) {
			this._textarea = new leap_bind_textarea(this, text_area);
		};
	}
} catch (e) {
}

/*--------------------------------------------------------------------------------------------------
 */

})();
