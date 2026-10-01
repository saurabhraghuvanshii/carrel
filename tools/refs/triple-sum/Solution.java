import java.util.*;

class Solution {
    public List<List<Integer>> tripleSum(int[] nums, int target) {
        int[] s = nums.clone();
        Arrays.sort(s);
        List<List<Integer>> out = new ArrayList<>();
        for (int i = 0; i + 2 < s.length; i++) {
            if (i > 0 && s[i] == s[i - 1]) {
                continue;
            }
            int lo = i + 1;
            int hi = s.length - 1;
            while (lo < hi) {
                int sum = s[i] + s[lo] + s[hi];
                if (sum < target) {
                    lo++;
                } else if (sum > target) {
                    hi--;
                } else {
                    out.add(List.of(s[i], s[lo], s[hi]));
                    while (lo < hi && s[lo] == s[lo + 1]) {
                        lo++;
                    }
                    lo++;
                    hi--;
                }
            }
        }
        return out;
    }
}
