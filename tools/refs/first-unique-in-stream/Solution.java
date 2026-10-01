import java.util.*;

// Ported from the owner's Ques/FirstNonReapt.java (printNonRepeating); it
// collects the answers, with '#' where the original printed -1.
class Solution {
    public String firstUnique(String str) {
        StringBuilder out = new StringBuilder();
        int freq[] = new int[26];
        Queue<Character> q = new LinkedList<>();
        for (int i = 0; i < str.length(); i++) {
            char ch = str.charAt(i);
            q.add(ch);
            freq[ch - 'a']++;
            while (!q.isEmpty() && freq[q.peek() - 'a'] > 1) {
                q.remove();
            }
            if (q.isEmpty()) {
                out.append('#');
            } else {
                out.append(q.peek());
            }
        }
        return out.toString();
    }
}
