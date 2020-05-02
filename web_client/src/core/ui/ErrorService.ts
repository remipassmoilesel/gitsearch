import {UiMessageService} from "./UiMessageService";

export class ErrorService {

    constructor(private uiMessages: UiMessageService) {
    }

    public onError(err: Error): void {
        console.error(err);

        let message: string;
        if (err.message) {
            message = `Error: ${err.message}`
        } else {
            message = `Something went wrong !`
        }

        this.uiMessages.error(message)
    }

}