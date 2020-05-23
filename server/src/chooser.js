const prefix = "http://localhost:8080"

class Chooser {
    constructor() {
        this.state = new State()
        this.mainDiv = this.addToDOM()

        /* this.getData()
            .then(data => this.updateDOMData(data))
            .catch(err => console.error(err)) */
    }

    addToDOM() {
        const chooser = document.createElement("div")
        chooser.setAttribute("class", "chooser")

        const header = document.createElement("h1")
        header.innerText = this.state.getHeder()
        chooser.appendChild(header)

        document.getElementById("chooserContainer").appendChild(chooser)

        return chooser
    }

    updateDOMData(data) {
        for (let i = 0; i < data.lenght; i++) {
            let chat = document.createElement("div")
            let name = document.createElement("h2")
        }
    }

    getData() {
        return fetch(this.state.getDataURL())
            .then(respone => respone.json())
            .then(jsonData => (jsonData.result))
    }
}

class State {
    constructor() {
        this.allStates  = [
            "showUsers",
            "showNewest",
            "showBiggest"
        ]

        this.cur = this.allStates[0]
    }

    getHeder() {
        switch (this.cur) {
            case this.allStates[0]:
                return "Dine chater: "

            case this.allStates[1]:
                return "Nyeste chater: "

            case this.allStates[2]:
                return "Største chater: "
        }
    }

    getDataURL() {
        switch (this.cur) {
            case this.allStates[0]:
                return prefix + "/api/usersChats"

            case this.allStates[0]:
                return prefix + "/api/newestChats"

            case this.allStates[0]:
                return prefix + "/api/bigestChats"
        }
    }
}