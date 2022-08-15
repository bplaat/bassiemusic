export function formatDuration(totalSeconds) {
    totalSeconds = Math.floor(totalSeconds);
    const hours = Math.floor(totalSeconds / 3600);
    const minutes = Math.floor((totalSeconds % 3600) / 60);
    const seconds = totalSeconds % 60;
    if (totalSeconds > 3600) {
        return hours + (minutes < 10 ? '0' + minutes : minutes) + ':' +
            (seconds < 10 ? '0' + seconds : seconds);
    }
    return minutes + ':' + (seconds < 10 ? '0' + seconds : seconds);
}
