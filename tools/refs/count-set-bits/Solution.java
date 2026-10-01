// Ported from the owner's bit manipulation/Oddoreven.java (countSetBits).
// Fixed: the loop ran while n > 0, so every negative number gave 0. It now
// runs while n != 0 and shifts with >>>, which brings in zeros at the top.
class Solution {
    public int countSetBits(int n) {
        int count = 0;
        while (n != 0) {
            if ((n & 1) != 0) {
                count++;
            }
            n = n >>> 1;
        }
        return count;
    }
}
