const dayjs = require('dayjs');

exports.calculateMP4StartTime = async (event) => {
    console.log(event);
    const startTime = event[0].ChannelInput.StartTime;
    const fiveMinutesAgo = dayjs(startTime).subtract(5, 'minute').format('YYYY-MM-DDTHH:mm:ss.SSSZ');
    console.log(fiveMinutesAgo);
    return fiveMinutesAgo;
}
