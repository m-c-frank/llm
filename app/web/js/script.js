let allMessages = []; // Global array to store messages

document.getElementById('send-button').addEventListener('click', sendMessage);

function sendMessage() {
    console.log("attempting to send message");
    const messageInput = document.getElementById('message-input');
    const message = messageInput.value.trim();
    if (!message) return;

    // Add message to UI and store it
    addMessage('User', message);

    // Send message to server
    fetch('/api', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ model: "mistral:instruct", messages: allMessages })
    })
    .then(response => response.json())
    .then(data => {
        // Assuming the server sends back the assistant's response
        addMessage('Assistant', data.reply);
    })
    .catch(error => console.error('Error:', error));

    // Clear the input field
    messageInput.value = '';
}

function addMessage(sender, text) {
    const messagesDiv = document.getElementById('messages');
    const messageElement = document.createElement('div');
    messageElement.classList.add('message');
    messageElement.textContent = `${sender}: ${text}`;
    messagesDiv.appendChild(messageElement);
    messagesDiv.scrollTop = messagesDiv.scrollHeight; // Auto-scroll to the latest message

    // Store the message in the global array
    allMessages.push({ sender, text });
}

function displayAllMessages() {
    // Optional function to display all stored messages
    console.log("All Messages:", allMessages);
}

