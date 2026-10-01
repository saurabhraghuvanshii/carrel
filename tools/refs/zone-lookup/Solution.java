import java.util.*;

class Solution {
    public int[] lookup(int[][] ranges, int[] codes) {
        int[][] sorted = ranges.clone();
        Arrays.sort(sorted, (a, b) -> Integer.compare(a[0], b[0]));
        int[] out = new int[codes.length];
        for (int i = 0; i < codes.length; i++) {
            int lo = 0, hi = sorted.length; // ranges before lo start at or below the code
            while (lo < hi) {
                int mid = (lo + hi) / 2;
                if (sorted[mid][0] <= codes[i]) {
                    lo = mid + 1;
                } else {
                    hi = mid;
                }
            }
            out[i] = lo > 0 && codes[i] <= sorted[lo - 1][1] ? sorted[lo - 1][2] : -1;
        }
        return out;
    }
}
