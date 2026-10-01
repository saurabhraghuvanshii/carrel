class Solution {
    public String removeTwins(String text) {
        StringBuilder kept = new StringBuilder();
        for (char ch : text.toCharArray()) {
            int last = kept.length() - 1;
            if (last >= 0 && kept.charAt(last) == ch) {
                kept.deleteCharAt(last);
            } else {
                kept.append(ch);
            }
        }
        return kept.toString();
    }
}
