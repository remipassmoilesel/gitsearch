enum MessageSeverity {
    INFO = "gs-info",
    ERROR = "gs-error",
}

export class UiMessageService {

    public info(content: string): void {
        this.message(MessageSeverity.INFO, content);
    }

    public error(content: string): void {
        this.message(MessageSeverity.ERROR, content);
    }

    private message(severity: MessageSeverity, content: string): void {
        const elem = document.querySelector("#gs-message-banner");
        if (!elem) {
            console.error(new Error("Message container not found"));
            return;
        }
        elem.innerHTML = `<div class='${severity}'>${content}</div>`
    }

}