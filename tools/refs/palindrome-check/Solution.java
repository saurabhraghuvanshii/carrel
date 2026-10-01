class Solution {
    public boolean sameBothWays(String text) {
        int i = 0;
        int j = text.length() - 1;
        while (i < j) {
            char a = text.charAt(i);
            char b = text.charAt(j);
            if (!Character.isLetterOrDigit(a)) {
                i++;
            } else if (!Character.isLetterOrDigit(b)) {
                j--;
            } else if (Character.toLowerCase(a) != Character.toLowerCase(b)) {
                return false;
            } else {
                i++;
                j--;
            }
        }
        return true;
    }
}
