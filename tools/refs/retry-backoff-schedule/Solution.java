class Solution {
    public long attempts(long base, long cap, long deadline) {
        long time = 0, count = 1, wait = base;
        while (wait < cap) {
            if (time + wait > deadline) {
                return count;
            }
            time += wait;
            count++;
            wait *= 2;
        }
        return count + (deadline - time) / cap;
    }
}
