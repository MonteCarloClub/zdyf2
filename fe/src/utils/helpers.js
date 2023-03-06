/**
 * 节流
 * @reference https://stackoverflow.com/questions/27078285/simple-throttle-in-javascript
 * @param {*} callback 
 * @param {*} limit 
 * @returns 
 */
export function throttle(callback, limit) {
    let obj;
    var waiting = false;                         // Initially, we're not waiting
    return function () {                         // We return a throttled function
        if (!waiting) {                          // If we're not waiting
            obj = this                           // For bind this in setTimeout
            waiting = true;                      // Prevent future invocations
            setTimeout(function () {             // After a period of time
                callback.apply(obj, arguments);  // Execute users function
                waiting = false;                 // And allow future invocations
            }, limit);
        }
    }
}