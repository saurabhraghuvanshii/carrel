import java.util.*;

class Solution {
    public int[] mostFrequent(int[] nums, int k) {
        Map<Integer, Integer> count = new HashMap<>();
        for (int v : nums) {
            count.merge(v, 1, Integer::sum);
        }
        List<List<Integer>> buckets = new ArrayList<>();
        for (int i = 0; i <= nums.length; i++) {
            buckets.add(new ArrayList<>());
        }
        for (Map.Entry<Integer, Integer> e : count.entrySet()) {
            buckets.get(e.getValue()).add(e.getKey());
        }
        int[] out = new int[k];
        int filled = 0;
        for (int f = nums.length; f > 0 && filled < k; f--) {
            for (int v : buckets.get(f)) {
                if (filled < k) {
                    out[filled++] = v;
                }
            }
        }
        return out;
    }
}
