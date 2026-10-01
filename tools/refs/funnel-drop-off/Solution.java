import java.util.*;

class Solution {
    public int[] funnel(int[][] events, int steps) {
        Map<Integer, Integer> reached = new HashMap<>();
        int[] counts = new int[steps];
        for (int[] e : events) {
            if (e[1] == reached.getOrDefault(e[0], 0) + 1) {
                reached.put(e[0], e[1]);
                counts[e[1] - 1]++;
            }
        }
        return counts;
    }
}
