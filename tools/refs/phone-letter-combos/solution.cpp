#include <string>
#include <vector>
using namespace std;

vector<string> words(string& digits) {
    static const string keys[] = {"", "", "abc", "def", "ghi", "jkl", "mno", "pqrs", "tuv", "wxyz"};
    vector<string> out = {""};
    for (char d : digits) {
        vector<string> next;
        for (auto& w : out)
            for (char ch : keys[d - '0']) next.push_back(w + ch);
        out = next;
    }
    return out;
}
