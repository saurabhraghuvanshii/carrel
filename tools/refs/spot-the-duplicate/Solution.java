import java.util.*;

class Solution {
    public boolean hasDuplicate(int[] nums) {
        Set<Integer> seen = new HashSet<>();
        for (int v : nums) {
            if (!seen.add(v)) {
                return true;
            }
        }
        return false;
    }
}
