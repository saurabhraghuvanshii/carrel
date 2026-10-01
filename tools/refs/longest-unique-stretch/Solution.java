import java.util.*;

class Solution {
    public int longestUnique(String text) {
        Map<Character, Integer> last = new HashMap<>();
        int best = 0;
        int start = 0;
        for (int i = 0; i < text.length(); i++) {
            Integer j = last.get(text.charAt(i));
            if (j != null && j >= start) {
                start = j + 1;
            }
            last.put(text.charAt(i), i);
            best = Math.max(best, i - start + 1);
        }
        return best;
    }
}
