import java.util.*;

class Solution {
    public int fewestCoins(int[] coins, int amount) {
        int[] best = new int[amount + 1];
        Arrays.fill(best, Integer.MAX_VALUE);
        best[0] = 0;
        for (int v = 1; v <= amount; v++) {
            for (int c : coins) {
                if (c <= v && best[v - c] != Integer.MAX_VALUE) {
                    best[v] = Math.min(best[v], best[v - c] + 1);
                }
            }
        }
        return best[amount] == Integer.MAX_VALUE ? -1 : best[amount];
    }
}
