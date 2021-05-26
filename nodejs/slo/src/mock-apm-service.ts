export const getApmDataMock = async (
    start: number,
    step: number,
    minutesRange: number,
): Promise<{ values: [number, string][] }> => {
    const results: [number, string][] = [];
    const steps = (minutesRange * 60) / step;
    for (let i = 0; i <= steps; i++) {
        const time = start + step * i;
        const value = Math.floor(Math.random() * 100) + 500;
        results.push([time, value.toString()]);
    }
    return new Promise((resolve) => {
        setTimeout(
            () =>
                resolve({
                    values: results,
                }),
            1000,
        );
    });
};
