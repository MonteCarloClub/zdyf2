// 提供将数据下载到本地文件的功能
export const FileDownloader = {
    methods: {
        saveFile(fileName, _) {
            // https://stackoverflow.com/questions/60810249/how-to-download-xlsx-file-from-a-server-response-in-javascript/64680613#64680613
            const url = window.URL.createObjectURL(new Blob([_]));
            const link = document.createElement("a");
            link.href = url;
            link.setAttribute("download", fileName);
            document.body.appendChild(link);
            link.click();
        },
    }
}