// 格式化时间
export const TimeFormat = {
    methods: {
        formatTime(unix_timestamp) {
            let a = "-";
            try {
                // Create a new JavaScript Date object based on the timestamp
                // multiplied by 1000 so that the argument is in milliseconds, not seconds.
                a = new Date(unix_timestamp * 1000);
                var year = a.getFullYear();
                var month = a.getMonth() + 1;
                var date = a.getDate();
                var hour = a.getHours();
                var min = a.getMinutes();
                var sec = a.getSeconds();
                var time = `${year}-${month}-${date} ${hour}:${min}:${sec}`;
                return time;
            } catch (error) {
                console.log(error);
            }

            return a
        },
    }
}