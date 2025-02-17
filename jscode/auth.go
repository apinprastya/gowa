package jscode

const InjectAuthJS = `function() {
	window.AuthStore = {};
	window.AuthStore.AppState = window.require('WAWebSocketModel').Socket;
	window.AuthStore.Cmd = window.require('WAWebCmd').Cmd;
	window.AuthStore.Conn = window.require('WAWebConnModel').Conn;
	window.AuthStore.OfflineMessageHandler = window.require('WAWebOfflineHandler').OfflineMessageHandler;
	window.AuthStore.PairingCodeLinkUtils = window.require('WAWebAltDeviceLinkingApi');
	window.AuthStore.Base64Tools = window.require('WABase64');
	window.AuthStore.RegistrationUtils = {
		...window.require('WAWebCompanionRegClientUtils'),
		...window.require('WAWebAdvSignatureApi'),
		...window.require('WAWebUserPrefsInfoStore'),
		...window.require('WAWebSignalStoreApi'),
	};
}`

const NeedAuthJS = `async function () {
	let state = window.AuthStore.AppState.state;
	if (state === 'OPENING' || state === 'UNLAUNCHED' || state === 'PAIRING') {
		// wait till state changes
		await new Promise(r => {
			window.AuthStore.AppState.on('change:state', function waitTillInit(_AppState, state) {
				if (state !== 'OPENING' && state !== 'UNLAUNCHED' && state !== 'PAIRING') {
					window.AuthStore.AppState.off('change:state', waitTillInit);
					r();
				} 
			});
		}); 
	}
	state = window.AuthStore.AppState.state;
	return state == 'UNPAIRED' || state == 'UNPAIRED_IDLE';
}`

const QrJS = `async function() {
	const registrationInfo = await window.AuthStore.RegistrationUtils.waSignalStore.getRegistrationInfo();
	const noiseKeyPair = await window.AuthStore.RegistrationUtils.waNoiseInfo.get();
	const staticKeyB64 = window.AuthStore.Base64Tools.encodeB64(noiseKeyPair.staticKeyPair.pubKey);
	const identityKeyB64 = window.AuthStore.Base64Tools.encodeB64(registrationInfo.identityKeyPair.pubKey);
	const advSecretKey = await window.AuthStore.RegistrationUtils.getADVSecretKey();
	const platform =  window.AuthStore.RegistrationUtils.DEVICE_PLATFORM;
	const getQR = (ref) => ref + ',' + staticKeyB64 + ',' + identityKeyB64 + ',' + advSecretKey + ',' + platform;
	
	window.onQRChangedEvent(getQR(window.AuthStore.Conn.ref)); // initial qr
	window.AuthStore.Conn.on('change:ref', (_, ref) => { window.onQRChangedEvent(getQR(ref)); }); // future QR changes
}`

const RequestPairCodeJS = `async function(phoneNumber) {
	window.AuthStore.PairingCodeLinkUtils.setPairingType('ALT_DEVICE_LINKING');
	await window.AuthStore.PairingCodeLinkUtils.initializeAltDeviceLinking();
	return window.AuthStore.PairingCodeLinkUtils.startAltLinkingFlow(phoneNumber, true);
}`
