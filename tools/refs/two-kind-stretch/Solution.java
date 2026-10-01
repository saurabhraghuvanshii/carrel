import java.util.*;

class Solution {
    public int longestTwoKinds(int[] kinds) {
        Map<Integer, Integer> count = new HashMap<>();
        int best = 0;
        int start = 0;
        for (int i = 0; i < kinds.length; i++) {
            count.merge(kinds[i], 1, Integer::sum);
            while (count.size() > 2) {
                if (count.merge(kinds[start], -1, Integer::sum) == 0) {
                    count.remove(kinds[start]);
                }
                start++;
            }
            best = Math.max(best, i - start + 1);
        }
        return best;
    }
}
