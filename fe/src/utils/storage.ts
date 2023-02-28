const MOCK_TX_LIST = 'key'

/**
 * 
 * @returns 
 */
export function getTxList() {
    const s = localStorage.getItem(MOCK_TX_LIST);
    if (s) {
        return JSON.parse(s);
    }
    return [];
}