interface Window {
    go: {
        main: {
            App: {
                GetInitialImage(): Promise<string>;
                ProcessImage(data: string): Promise<any>;
                ProcessImageFile(path: string): Promise<any>;
                Greet(name: string): Promise<string>;
            }
        }
    }
}
