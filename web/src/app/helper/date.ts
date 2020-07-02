export function convertToDDMMYYY(timestamp: number) {
    const date = new Date(timestamp * 1000);
    let day: number | string = date.getDate();
    let month: number | string = date.getMonth();
    month = month + 1;
    if ((String(day)).length == 1)
        day = '0' + day;
    if ((String(month)).length == 1)
        month = '0' + month;

    return `${day}.${month}.${date.getFullYear()}`
}