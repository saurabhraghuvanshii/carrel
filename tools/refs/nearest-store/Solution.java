import java.util.*;

class Solution {
    public int[] nearest(int[] stores, int[] customers) {
        int[] sorted = stores.clone();
        Arrays.sort(sorted);
        int[] out = new int[customers.length];
        for (int i = 0; i < customers.length; i++) {
            int lo = 0, hi = sorted.length; // first store at or after the customer
            while (lo < hi) {
                int mid = (lo + hi) / 2;
                if (sorted[mid] < customers[i]) {
                    lo = mid + 1;
                } else {
                    hi = mid;
                }
            }
            int best = Integer.MAX_VALUE;
            if (lo < sorted.length) {
                best = sorted[lo] - customers[i];
            }
            if (lo > 0) {
                best = Math.min(best, customers[i] - sorted[lo - 1]);
            }
            out[i] = best;
        }
        return out;
    }
}
