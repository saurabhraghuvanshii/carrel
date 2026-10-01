import java.util.*;

class Solution {
    public boolean containsRearrangement(String pattern, String text) {
        int m = pattern.length();
        if (m > text.length()) {
            return false;
        }
        int[] need = new int[26];
        int[] have = new int[26];
        for (int i = 0; i < m; i++) {
            need[pattern.charAt(i) - 'a']++;
            have[text.charAt(i) - 'a']++;
        }
        for (int i = m; ; i++) {
            if (Arrays.equals(need, have)) {
                return true;
            }
            if (i == text.length()) {
                return false;
            }
            have[text.charAt(i) - 'a']++;
            have[text.charAt(i - m) - 'a']--;
        }
    }
}
