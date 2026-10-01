import java.util.*;

class Solution {
    public int countBatches(int[][] requests, int limit, int window) {
        Map<Long, int[]> open = new HashMap<>();
        int count = 0;
        for (int[] r : requests) {
            long zone = (long) Math.floorDiv(r[1], 100) * 1_000_003L + Math.floorDiv(r[2], 100);
            int[] b = open.get(zone);
            if (b != null && r[0] - b[0] <= window && b[1] < limit) {
                b[1]++;
            } else {
                open.put(zone, new int[] {r[0], 1});
                count++;
            }
        }
        return count;
    }
}
