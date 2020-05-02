export class Logger {

    public static error(message: string, data?: any): void {
        console.error(message)
        if(data){
            console.error(data)
        }
    }

}