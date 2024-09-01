// worker.onmessage = function (e) {
//   console.log(new Date());
//   delay.push({
//     time: new Date(),
//     forward: Math.round(Math.random() * 100),
//     gateway: Math.round(100 + Math.random() * 50),
//     forward_loc: '北京',
//     gateway_loc: '莫斯科',
//   })
// }

// worker = new Worker('worker.js');
// worker.onmessage = function (e) { console.log(new Date()); }
// worker.postMessage('start')

var interval
self.onmessage = function (e) {
  if (e.data > 0) {
    if (interval == null) { clearInterval(interval) }
    interval = setInterval(() => {
      self.postMessage(`signal`)
    }, e.data)
  } else if (interval !== null && e.data == 0) {
    clearInterval(interval)
  }
};