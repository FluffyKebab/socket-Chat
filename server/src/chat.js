const defultMessage = "Du har ikke fått noen melinger på denne chaten"

class Chat {
    constructor(chatName) {
        this.chatName = chatName
        this.chatDiv = this.addChatToDOM()
        this.ws = new WebSocket(`ws://localhost:8080/ws/messaageConn?name=${this.chatName}`)
        
        this.ws.addEventListener("message", e => {
            if (this.chatDiv.innerText == defultMessage) {
                this.chatDiv.innerText = ""
            }

            const message = JSON.parse(e.data)
            message.isMessage = true

            this.addMessageToDOM(message)
            this.addMessageToStorage(message)
            
        })

        this.logDate(false)
        this.updateScroll()
        this.addStorage()
    }

    addChatToDOM() {
        const newChat = document.createElement("div")
        newChat.setAttribute("class", "chat")

        const name = document.createElement("h1")
        name.innerText = this.chatName
        newChat.appendChild(name)

        const messageDiv = document.createElement("div")
        messageDiv.setAttribute("class", "messages")
        newChat.appendChild(messageDiv)

        const form = document.createElement("form")
        const inputText = document.createElement("input")
        const submit = document.createElement("input")

        inputText.setAttribute("type", "text")
        inputText.setAttribute("name", "text")
        inputText.setAttribute("autocomplete", "off")
        
        submit.setAttribute("type", "submit")
        submit.setAttribute("value", "Send")
        
        form.appendChild(inputText)
        form.appendChild(submit)

        newChat.appendChild(form)

        form.addEventListener("submit", e => {
            e.preventDefault()

            if (this.ws.readyState === WebSocket.OPEN && inputText.value != "") {
                this.ws.send(inputText.value)
                inputText.value = ""
            }
        })

        document.getElementById("chatContainer").appendChild(newChat)

        return messageDiv
    }

    addStorage() {
        const storageValue = JSON.parse(localStorage.getItem(this.chatName))

        if (storageValue == null || storageValue.length < 1) {
            localStorage.setItem(this.chatName, JSON.stringify([]))
            this.chatDiv.innerText = defultMessage
            return
        }

        for (let i = 0; i < storageValue.length; i++) {
            this.addMessageToDOM(storageValue[i])
        }
    }

    addMessageToDOM(messageData) {
        const messageDiv = document.createElement("div")

        if (!messageData.isMessage) {
            const date = document.createElement("h4")
            date.innerText =  messageData.date
            messageDiv.appendChild(date)
        }

        if (messageData.isMessage) {
            messageDiv.classList.add("message")
            if (messageData.isUsersMessage) {
                messageDiv.classList.add("messageYour")
            }

            const name = document.createElement("h4")
            name.innerText = messageData.displayName
            messageDiv.appendChild(name)

            const p = document.createElement("p")
            p.innerText = messageData.body
            messageDiv.appendChild(p)
        }   

        this.chatDiv.appendChild(messageDiv)
        
    }

    addMessageToStorage(messageData) {
        let storage = JSON.parse(localStorage.getItem(this.chatName))

        if (!storage) {
            storage = []
        }

        storage.push(messageData)
        localStorage.setItem(this.chatName, JSON.stringify(storage))
    }

    updateScroll() {
        let scrolled = false

        setInterval(() => {
            if (!scrolled) {
                this.chatDiv.scrollTop = this.chatDiv.scrollHeight
            }
        }, 100)

        this.chatDiv.scrollTop = this.chatDiv.scrollHeigth

        this.chatDiv.addEventListener("scroll", () => {
            scrolled = !((this.chatDiv.scrollHeight - this.chatDiv.offsetHeight) <= this.chatDiv.scrollTop)
        })
    }

    logDate(isLogof) {
        const dateNow = new Date()
        const logofInfo = {
            isMessage: false,
            date: `Du logget ${isLogof ? "av" : "på"}  ${dateNow.getDate()}.${dateNow.getMonth() + 1}.${dateNow.getFullYear()} klokken ${ dateNow.getHours() }:${ dateNow.getMinutes() }`,
        }

        this.addMessageToStorage(logofInfo)
        this.addMessageToDOM(logofInfo)
    }
}