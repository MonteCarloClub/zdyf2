function pad2(n: number) {
  return (n < 10 ? "0" : "") + n;
}

export function getTime() {
  var myDate = new Date();
  const y = myDate.getFullYear(); //获取完整的年份(4位,1970-????)
  const m = 1 + myDate.getMonth(); //获取当前月份(0-11,0代表1月)
  const d = myDate.getDate(); //获取当前日(1-31)

  const h = myDate.getHours(); //获取当前小时数(0-23)
  const mi = myDate.getMinutes(); //获取当前分钟数(0-59)
  const s = myDate.getSeconds(); //获取当前秒数(0-59)
  return `${y}-${pad2(m)}-${pad2(d)} ${pad2(h)}:${pad2(mi)}:${pad2(s)} `;
}
