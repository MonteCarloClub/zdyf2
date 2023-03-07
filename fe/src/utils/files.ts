/**
 * 将内容下载为本地文件
 * @param content 文本内容
 * @param fileName 带后缀的文件名
 */
export function download(content: string, fileName: string = "default.txt") {
  var a = document.createElement("a");
  var file = new Blob([content], { type: "text/plain" });
  a.href = URL.createObjectURL(file);
  a.download = fileName;
  a.click();
}

type ReadfileCallback = (file: File, content: string) => void;
type onError = (error: string) => void;
/**
 * 读取本地文件内容
 * @param cb 回调
 */
export function readLocalFile(cb: ReadfileCallback, error: onError) {
  var input = document.createElement("input");
  input.type = "file";

  input.onchange = (e: Event) => {
    const input = e.target as HTMLInputElement;
    
    if (input.files === null || input.files.length === 0) {
      error('files is null')
      return;
    }
    const file = input.files[0];
    // this.certFileName = file.name;
    // setting up the reader
    var reader = new FileReader();
    reader.readAsText(file, "UTF-8");

    // here we tell the reader what to do when it's done reading...
    reader.onload = (readerEvent) => {
      if (readerEvent.target === null) {
        error('readerEvent.target is null')
        return;
      }
      var content = readerEvent.target.result; // this is the content!
      cb(file, content as string);
    };
  };

  input.click();
}