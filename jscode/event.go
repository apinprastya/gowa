package jscode

const StateChangeJS = `function() {
	window.AuthStore.AppState.on('change:state', (_AppState, state) => { window.onAuthAppStateChangedEvent(state); });
	window.AuthStore.AppState.on('change:hasSynced', () => { window.onAppStateHasSyncedEvent(); });
	window.AuthStore.Cmd.on('logout', async () => {
		await window.onLogoutEvent();
	});
}`

const MessageEventJS = `() => {
	window.Store.Msg.on('change:ack', (msg, ack) => {
		window.onMessageAckEvent(window.WWebJS.getMessageModel(msg).id.id, ack); 
	});
}`
