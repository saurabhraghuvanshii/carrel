import java.util.*;

class Solution {
    public long minTotalDistance(int[] orders, int[] couriers) {
        int[] o = orders.clone(), c = couriers.clone();
        Arrays.sort(o);
        Arrays.sort(c);
        int n = o.length, m = c.length;
        long[] prev = new long[m + 1];
        for (int i = 1; i <= n; i++) {
            long[] cur = new long[m + 1];
            Arrays.fill(cur, Long.MAX_VALUE / 2);
            for (int j = i; j <= m; j++) {
                cur[j] = Math.min(cur[j - 1], prev[j - 1] + Math.abs(o[i - 1] - c[j - 1]));
            }
            prev = cur;
        }
        return prev[m];
    }
}
