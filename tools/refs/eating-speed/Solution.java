class Solution {
    public int slowestSpeed(int[] jobs, int h) {
        int lo = 1;
        int hi = 1;
        for (int j : jobs) {
            hi = Math.max(hi, j);
        }
        while (lo < hi) {
            int mid = lo + (hi - lo) / 2;
            long hours = 0;
            for (int j : jobs) {
                hours += (j + (long) mid - 1) / mid;
            }
            if (hours <= h) {
                hi = mid;
            } else {
                lo = mid + 1;
            }
        }
        return lo;
    }
}
