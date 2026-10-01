class Solution {
    public String reverseWords(String text) {
        StringBuilder out = new StringBuilder();
        int end = text.length();
        for (int i = text.length() - 1; i >= -1; i--) {
            if (i == -1 || text.charAt(i) == ' ') {
                if (i + 1 < end) {
                    if (out.length() > 0) {
                        out.append(' ');
                    }
                    out.append(text, i + 1, end);
                }
                end = i;
            }
        }
        return out.toString();
    }
}
