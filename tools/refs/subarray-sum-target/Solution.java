import java.util.*;

class Solution {
    public int countSubarrays(int[] nums, int target) {
        Map<Integer, Integer> seen = new HashMap<>();
        seen.put(0, 1);
        int sum = 0;
        int count = 0;
        for (int v : nums) {
            sum += v;
            count += seen.getOrDefault(sum - target, 0);
            seen.merge(sum, 1, Integer::sum);
        }
        return count;
    }
}
