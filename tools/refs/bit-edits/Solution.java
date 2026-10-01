// Ported from the owner's bit manipulation/Oddoreven.java (getIthBit,
// setIthbit, clearIthbit, updateIthbit). Fixed: clearIthbit used ~1 << i,
// which is (~1) << i and clears bits 0 to i, not just bit i; it is now
// ~(1 << i).
class Solution {
    public int getBit(int n, int i) {
        int bitmask = 1 << i;
        if ((n & bitmask) == 0) {
            return 0;
        } else {
            return 1;
        }
    }

    public int setBit(int n, int i) {
        int bitmask = 1 << i;
        return n | bitmask;
    }

    public int clearBit(int n, int i) {
        int bitmask = ~(1 << i);
        return n & bitmask;
    }

    public int updateBit(int n, int i, int newBit) {
        n = clearBit(n, i);
        int bitMask = newBit << i;
        return n | bitMask;
    }
}
