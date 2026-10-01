import java.util.*;

class Solution {
    private static final String[] KEYS = {"", "", "abc", "def", "ghi", "jkl", "mno", "pqrs", "tuv", "wxyz"};

    public List<String> words(String digits) {
        List<String> out = new ArrayList<>();
        walk(digits, 0, new StringBuilder(), out);
        return out;
    }

    private void walk(String digits, int i, StringBuilder cur, List<String> out) {
        if (i == digits.length()) {
            out.add(cur.toString());
            return;
        }
        for (char ch : KEYS[digits.charAt(i) - '0'].toCharArray()) {
            cur.append(ch);
            walk(digits, i + 1, cur, out);
            cur.deleteCharAt(cur.length() - 1);
        }
    }
}
