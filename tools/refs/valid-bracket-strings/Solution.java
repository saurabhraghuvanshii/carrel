import java.util.*;

class Solution {
    public List<String> balancedStrings(int n) {
        List<String> out = new ArrayList<>();
        walk(n, 0, 0, new StringBuilder(), out);
        return out;
    }

    private void walk(int n, int open, int close, StringBuilder cur, List<String> out) {
        if (cur.length() == 2 * n) {
            out.add(cur.toString());
            return;
        }
        if (open < n) {
            cur.append('(');
            walk(n, open + 1, close, cur, out);
            cur.deleteCharAt(cur.length() - 1);
        }
        if (close < open) {
            cur.append(')');
            walk(n, open, close + 1, cur, out);
            cur.deleteCharAt(cur.length() - 1);
        }
    }
}
