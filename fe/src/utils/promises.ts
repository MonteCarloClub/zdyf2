/**
 * 限制最大并发数的 Promise.all
 * @param tasks 一组异步任务
 * @param limit 同时执行的任务数量上线
 * @returns 执行结果
 */
export function promiseLimit<R>(tasks: Promise<R>[], limit: number): Promise<R[]> {
  return new Promise((resolve) => {
    const results: R[] = [];
    const queue = [...tasks];
    let index = 0;
    let resolvedCount = 0;

    function next(i: number) {
      const task = queue.shift();
      if (!task || i === tasks.length) {
        if (resolvedCount === tasks.length) {
          resolve(results);
        }
        return;
      }

      Promise.resolve(task).then((result) => {
        results[i] = result;
        resolvedCount++;
        next(index);
      });

      index++;
    }

    for (let i = 0; i < limit && queue.length > 0; i++) {
      next(index);
    }
  });
}

/**
 * 限制最大并发数执行一组任务
 * @param tasks 一组异步任务
 * @param limit 同时执行的任务数量上线
 * @param callback 每个任务 resolve 时的回调
 */
export function parallelWithLimit<R>(
  tasks: Promise<R>[],
  limit: number,
  callback: (index: number, res: R) => void
) {
  const queue = [...tasks];
  let index = 0;
  let resolvedCount = 0;

  function next(i: number) {
    const task = queue.shift();
    if (!task || i === tasks.length) {
      return;
    }
    
    task.then((result) => {
      callback(i, result);
      resolvedCount++;
      next(index);
    });

    index++;
  }

  for (let i = 0; i < limit && queue.length > 0; i++) {
    next(index);
  }
}
