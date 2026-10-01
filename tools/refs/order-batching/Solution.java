class Solution {
    public int smallestCapacity(int[] weights, int k) {
        int lo = 0, hi = 0;
        for (int w : weights) {
            lo = Math.max(lo, w);
            hi += w;
        }
        while (lo < hi) {
            int mid = lo + (hi - lo) / 2;
            if (vans(weights, mid) <= k) {
                hi = mid;
            } else {
                lo = mid + 1;
            }
        }
        return lo;
    }

    private int vans(int[] weights, int capacity) {
        int count = 1, load = 0;
        for (int w : weights) {
            if (load + w > capacity) {
                count++;
                load = 0;
            }
            load += w;
        }
        return count;
    }
}
