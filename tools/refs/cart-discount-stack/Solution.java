import java.util.*;

class Solution {
    public long finalTotal(int[] prices, int[][] coupons) {
        ArrayDeque<int[]> active = new ArrayDeque<>();
        for (int[] c : coupons) {
            if (c[0] != 0) {
                active.addLast(c);
            } else if (!active.isEmpty()) {
                active.pollLast();
            }
        }
        long total = 0;
        for (int p : prices) {
            total += p;
        }
        for (int[] c : active) {
            if (c[0] == 1) {
                total -= total * c[1] / 100;
            } else {
                total = Math.max(0, total - c[1]);
            }
        }
        return total;
    }
}
