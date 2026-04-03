package jscode

const SendMessageJS = `async (message) => {
	const chat = await window.WWebJS.getChat(message.chatId, {getAsModel: false});

	if (!chat) return null;

	const msg = await window.WWebJS.sendMessage(
                    chat,
                    message.message, 
					message.options,
                );
				
	return msg ? window.WWebJS.getMessageModel(msg) : undefined;
};`
