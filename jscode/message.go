package jscode

const SendMessageJS = `async (message) => {
	const chatWid = window.Store.WidFactory.createWid(message.chatId);
	const chat = await window.Store.Chat.find(chatWid);

	const msg = await window.WWebJS.sendMessage(chat, message.message, message.options, false);
	return window.WWebJS.getMessageModel(msg).id.id;
};`
