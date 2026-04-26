package jscode

const StateChangeJS = `function() {
	window.require('WAWebSocketModel')
                .Socket.on('change:state', (_AppState, state) => {
                    window.onAuthAppStateChangedEvent(state);
                });
	window.require('WAWebSocketModel')
                .Socket.on('change:hasSynced', () => {
                    window.onAppStateHasSyncedEvent();
                });
	const Cmd = window.require('WAWebCmd').Cmd;
	Cmd.on('logout', async () => {
		await window.onLogoutEvent();
	});
	Cmd.on('logout_from_bridge', async () => {
		await window.onLogoutEvent();
	});
}`

const MessageEventJS = `() => {
	const { Msg } = window.require('WAWebCollections');
	Msg.on('change:ack', (msg, ack) => {
		window.onMessageAckEvent(window.WWebJS.getMessageModel(msg).id.id, ack); 
	});
}`
