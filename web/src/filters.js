export function formatDuration(totalSeconds) {
    totalSeconds = Math.round(totalSeconds);
    const hours = Math.floor(totalSeconds / 3600);
    const minutes = Math.floor((totalSeconds % 3600) / 60);
    const seconds = totalSeconds % 60;
    if (totalSeconds > 3600) {
        return hours + (minutes < 10 ? '0' + minutes : minutes) + ':' + (seconds < 10 ? '0' + seconds : seconds);
    }
    return minutes + ':' + (seconds < 10 ? '0' + seconds : seconds);
}

export function formatBytes(bytes) {
    if (bytes > 1e12) return `${(bytes / 1e12).toFixed(2)} TB`;
    if (bytes > 1e9) return `${(bytes / 1e9).toFixed(2)} GB`;
    if (bytes > 1e6) return `${(bytes / 1e6).toFixed(2)} MB`;
    if (bytes > 1e3) return `${Math.ceil(bytes / 1e3)} KB`;
    return `${bytes} B`;
}
