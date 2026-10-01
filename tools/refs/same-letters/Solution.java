class Solution {
    public boolean sameLetters(String first, String second) {
        if (first.length() != second.length()) {
            return false;
        }
        int[] count = new int[26];
        for (int i = 0; i < first.length(); i++) {
            count[first.charAt(i) - 'a']++;
            count[second.charAt(i) - 'a']--;
        }
        for (int c : count) {
            if (c != 0) {
                return false;
            }
        }
        return true;
    }
}
