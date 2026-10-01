#include <cctype>
#include <string>
using namespace std;

bool sameBothWays(string& text) {
    int i = 0, j = (int)text.size() - 1;
    while (i < j) {
        unsigned char a = text[i], b = text[j];
        if (!isalnum(a)) i++;
        else if (!isalnum(b)) j--;
        else if (tolower(a) != tolower(b)) return false;
        else {
            i++;
            j--;
        }
    }
    return true;
}
