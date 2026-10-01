import java.util.*;

class Solution {
    public int longestRun(int[] nums) {
        Set<Integer> set = new HashSet<>();
        for (int v : nums) {
            set.add(v);
        }
        int best = 0;
        for (int v : set) {
            if (set.contains(v - 1)) {
                continue;
            }
            int length = 1;
            while (set.contains(v + length)) {
                length++;
            }
            best = Math.max(best, length);
        }
        return best;
    }
}
